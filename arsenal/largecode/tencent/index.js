'use strict';
function cpuIntensiveCalculation(baseNumber) {
  var iterationCount = 50000 * Math.pow(baseNumber, 3);
  var result = 0;
  for (var i = iterationCount; i >= 0; i--) {
    result += Math.atan(i) * Math.tan(i);
  }
}
function getDuration(startTime) {
  var end = process.hrtime(startTime);
  return end[1] + (end[0] * 1e9);
}


function runTest(sleep_time) {
  cpuIntensiveCalculation(sleep_time);
}

function isWarm() {
  var is_warm = process.env.warm ? true : false;
  process.env.warm = true;
  return is_warm;
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
