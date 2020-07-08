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
var wait = ms => new Promise((r, j) => setTimeout(r, ms));

function getDuration(startTime) {
    var end = process.hrtime(startTime);
    return end[1] + (end[0] * 1e9);
}


async function runTest(sleep_time){
    await wait(sleep_time);
}

function isWarm() {
    var is_warm = process.env.warm ? true : false;
    process.env.warm = true;
    return is_warm;
}


exports.handler = async (req, resp, context) => {
    var startTime = process.hrtime();
    var params = {
        path: req.path,
        queries: req.queries,
        headers: req.headers,
        method : req.method,
        requestURI : req.url,
        clientIP : req.clientIP,
    }
    await runTest(params.queries['sleep']);
    var reused = isWarm();
    var duration = getDuration(startTime);
    // resp.send(JSON.stringify(params, null, '    '));
    resp.send(JSON.stringify({
        reused: reused,
        duration: duration
    }));

    /*
    getFormBody(req, function(err, formBody) {
        for (var key in req.queries) {
          var value = req.queries[key];
          resp.setHeader(key, value);
        }
        params.body = formBody;
        console.log(formBody);
        resp.send(JSON.stringify(params));
    });
    */
}