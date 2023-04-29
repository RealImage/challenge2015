const Joi = require('joi');
const UtilService = require('../utils');

module.exports = {
    degreeOfSeparation: (req, res, next) => {
        try {
            const schema = Joi.object({
                actor1: Joi.string().required().min(1),
                actor2: Joi.string().required().min(1),
            }).validate(req.query || {});

            UtilService.validateJoi(schema, req, res, next);

        } catch (err) {
            return UtilService.sendResponse(res, 500, 'Something Went Wrong', {});

        }
    }
}