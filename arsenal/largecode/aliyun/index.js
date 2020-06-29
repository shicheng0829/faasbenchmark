var getRawBody = require('raw-body');
var getFormBody = require('body/form');
var body = require('body');


/*
To enable the initializer feature (https://help.aliyun.com/document_detail/156876.html)
please implement the initializer function as below：
exports.initializer = (context, callback) => {
  console.log('initializing');
  callback(null, '');
};
*/
function cpuIntensiveCalculation(baseNumber) {
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
  cpuIntensiveCalculation(intensityLevel);
}



exports.handler = async (req, resp, context) => {
  var startTime = process.hrtime();
  // console.log(req.queries['level'])
  runTest(req.queries['level']);
  var reused = isWarm();
  var duration = getDuration(startTime);
  resp.send(JSON.stringify({
    reused: reused,
    duration: duration
  }));

}