const constant = require("../constant");
const { default: axios } = require("axios");

module.exports = {
    sendResponse: (res, statusCode, message, data = {}) => {
        res.status(statusCode).json({
            message: typeof message === "string" ? message : "",
            data: data,
        });
    },

    validateJoi: (schema, req, res, next) => {
        if (schema.error) {
            const errMsg = schema.error.details[0].message;
            return res.status(400).json({
                message: errMsg,
                data: {}
            });
        } else {
            req.value = schema.value;
            return next();
        }
    },
    fetchDataFromUrl: function (url) {
        return axios(constant.UrlPrefix + url).then((data) => {
            if (data.status == 200)
                return data.data
        }).catch((err) => {
            if (err.response.status == 403) {
                return {
                    error: true,
                    url: url
                }
            };
            throw 'INVALID_URL'
        })
    },


    pushDataInMap: (map, moviesArray, actorName) => {
        moviesArray.forEach(element => {
            if (map.has(element.url) && map.get(element.url) === actorName) {
                map.delete(element.url);
            }
            else
                map.set(element.url, actorName)
        });
    },

}