'use strict';
const MEGABYTE = 1024 * 1024;

function memIntensive(level) {
  var available_memory = 512;
  // console.log(env.AWS_LAMBDA_FUNCTION_MEMORY_SIZE)
  let amountInMB = available_memory - (available_memory / 10) * (4 - level);
  // console.log(parseInt(amountInMB));
  Buffer.alloc(amountInMB * MEGABYTE, 'a');
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


function runTest(intensityLevel) {
  memIntensive(intensityLevel)
}

exports.main_handler = (event, context, callback) => {
  var startTime = process.hrtime();
  runTest(event["queryString"]["level"]);

  var reused = isWarm();
  var duration = getDuration(startTime);

  return {
    "reused": reused,
    "duration": duration
  }
};
