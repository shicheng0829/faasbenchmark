'use strict';
const PATH = '/tmp/faastest';
const proc = require('child_process');

function ioIntensive(baseNumber) {
  var amountInMB = 20 ** (baseNumber - 1);
  var out = proc.spawnSync('dd', ['if=/dev/zero', `of=${PATH}`, `bs=${amountInMB}M`, 'count=1', 'oflag=direct']);
  if (out.status !== 0)
    return out.stderr.toString();
}


function isWarm() {
  var is_warm = process.env.warm ? true : false;
  process.env.warm = true;
  return is_warm;
}

function getDuration(startTime) {
  var end = process.hrtime(startTime);
  return end[1] + (end[0] * 1e9);
}

function getLevel(event) {
  let intensityLevel = event.level ? parseInt(event.level) : null;
  if (!intensityLevel || intensityLevel < 1) {
    return {"error": "invalid level parameter"};
  }
  return intensityLevel;
}

function getParameters(event) {
  return getLevel(event);
}

function runTest(intensityLevel){
  ioIntensive(intensityLevel);
}


exports.main_handler = async (event, context, callback) => {
  var startTime = process.hrtime();
  runTest(event["queryString"]["level"]);
  var reused = isWarm();
  var duration = getDuration(startTime);
  return {
    "reused": reused,
    "duration": duration
  }
};