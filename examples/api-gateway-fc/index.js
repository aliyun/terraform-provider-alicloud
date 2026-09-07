'use strict';

var fs = require("fs")

exports.handler = function (event, context, callback) {

    console.log("request: " + JSON.stringify(event.toString()));
    // convert event to JSON object
    event = JSON.parse(event.toString());
    var query = event.queryParameters || {};
    var showResponse = query.response || 'html';

    switch (showResponse) {
    case "html":
        var htmlResponse = {
            isBase64Encoded: true,
            statusCode: 200,
            headers: {
                "Content-type": "text/html; charset=utf-8"
            },
            // base64 encode body so it can be safely returned as JSON value
            body: new Buffer("<html><h1>hello FunctionCompute</h1></html>").toString('base64')
        }
        callback(null, htmlResponse);
        break;
    case "json":
        var jsonResponse = {
            isBase64Encoded: true,
            statusCode: 200,
            headers: {
                "Content-type": "application/json"
            },
            // base64 encode body so it can be safely returned as JSON value
            body: new Buffer('{"hello": "FunctionCompute"}').toString('base64')
        }
        callback(null, jsonResponse);
        break;
    case "image":
        fs.readFile(process.env.FC_FUNC_CODE_PATH + "/fc.png", (err, data) => {
            if (err) {
                callback(null, err);
                return;
            }
            var imageResponse = {
                isBase64Encoded: true,
                statusCode: 200,
                headers: {
                    "Content-type": "image/png"
                },
                // base64 encode body so it can be safely returned as JSON value
                body: data.toString('base64')
            }
            callback(null, imageResponse);
            return;
        });
        break;
    default:
        var htmlResponse = {
            isBase64Encoded: true,
            statusCode: 200,
            headers: {
                "Content-type": "text/html; charset=utf-8"
            },
            // base64 encode body so it can be safely returned as JSON value
            body: new Buffer("<html>Please hit the url with response=html, reponse=json, or reponse=image</html>").toString('base64')
        }
        callback(null, htmlResponse);
    }
};