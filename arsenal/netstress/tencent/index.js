'use strict';
const http = require('http');
const fs = require('fs');

const files = {1: '/files/1Mb.dat', 2: '/files/10Mb.dat', 3: '/files/100Mb.dat'};

async function networkIntensive(level) {
  // console.log("start download")
  const writable = fs.createWriteStream('/dev/null');
  await new Promise((resolve) => http.get({
    host: `www.ovh.net`,
    port: 80,
    path: files[level]
  }, (res) => {
    var download = res.pipe(writable);
    download.on('close', () => resolve(res));
  }));
  // console.log("finish download")
}

function getDuration(startTime) {
  var end = process.hrtime(startTime);
  return end[1] + (end[0] * 1e9);
}


async function runTest(level) {
  await networkIntensive(level);
}

function isWarm() {
  var is_warm = process.env.warm ? true : false;
  process.env.warm = true;
  return is_warm;
}

exports.main_handler = async (event, context, callback) => {
  var startTime = process.hrtime();
  await runTest(event["queryString"]["level"]);

  var reused = isWarm();
  var duration = getDuration(startTime);

  return {
    "reused": reused,
    "duration": duration
  }
};
