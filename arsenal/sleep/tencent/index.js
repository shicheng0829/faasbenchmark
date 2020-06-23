'use strict';
var wait = ms => new Promise((r, j) => setTimeout(r, ms));

function getDuration(startTime) {
  var end = process.hrtime(startTime);
  return end[1] + (end[0] * 1e9);
}

function getSleep(event) {
  let sleep_time = event["queryString"]["sleep"] ? parseInt(event.sleep) : null;
  if (!sleep_time && sleep_time !== 0) {
    return { "error": event };
  }
  return sleep_time;
}

function getParameters(event) {
  return getSleep(event);
}

async function runTest(sleep_time) {
  await wait(sleep_time);
}

function isWarm() {
  var is_warm = process.env.warm ? true : false;
  process.env.warm = true;
  return is_warm;
}

exports.main_handler = async (event, context, callback) => {
  var startTime = process.hrtime();
  // return event["queryString"]["sleep"]
  let params = getParameters(event);
  // return { "body": params["queryString"]}
  // if (params.error) {
  //   return { "body": `{"error": ${params.error}}` }
  // }

  await runTest(params);

  var reused = isWarm();
  var duration = getDuration(startTime);

  return {
    "reused": reused,
    "duration": duration
  }
};
