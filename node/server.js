const express = require('express');
const { chain } = require('stream-chain');
const { parser } = require('stream-json');
const { streamArray } = require('stream-json/streamers/StreamArray');
const { Writable, Readable } = require('stream');
const cluster = require('cluster');
const os = require('os');
const process = require('process');

function main() {
  const app = express();

  app.get('/status', function (req, res) {
    res.send('Ok');
  });

  app.post('/filter', filter);

  app.listen(8080);
}

function filter(req, res) {
  const filteredStream = new Readable({
    read() {},
  });

  const pipeline = chain([req, parser(), streamArray()]);

  const processStream = new Writable({
    objectMode: true,
    write({ value }, encoding, callback) {
      if (value.version > 5) {
        filteredStream.push(JSON.stringify(value) + '\n');
      }
      callback();
    },
  });

  processStream.on('finish', () => {
    filteredStream.push(null);
  });

  pipeline.on('error', (err) => {
    console.error('Error while processing JSON:', err);
    res.status(500).send('Error processing JSON.');
  });

  processStream.on('error', (err) => {
    console.error('Error in processing stream:', err);
    res.status(500).send('Error processing JSON.');
  });

  pipeline.pipe(processStream);

  filteredStream.pipe(res);
}

if (cluster.isPrimary) {
  console.log(`Primary ${process.pid} is running`);

  for (let i = 0; i < os.availableParallelism(); i++) {
    cluster.fork();
  }

  cluster.on('exit', (worker, code, signal) => {
    console.log(`worker ${worker.process.pid} died`);
  });
} else {
  main();
  console.log(`Worker ${process.pid} started`);
}
