const express = require('express');
const validation = require('./validation');
const controller = require('./controller');
const app = express()

app.get('/degree-of-separation', validation.degreeOfSeparation, controller.degreeOfSeparation)


app.listen(4000, () => {
    console.log("App is running on the port>>>>4000")
});
