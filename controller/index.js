
const utils = require("../utils")

module.exports = {
    degreeOfSeparation: async (req, res, next) => {
        try {
            let separatedMovies = new Map();
            let includingMovies = [];
            let moviesCheck = [];
            let moviesData = await Promise.all([utils.fetchDataFromUrl(req.value.actor1), utils.fetchDataFromUrl(req.value.actor2)])
                .catch((err) => {
                    throw 'ACTOR_MOVIES_URL_ISSUE'
                })

            if (!moviesData[0]?.movies || !moviesData[1]?.movies)
                throw 'MOVIES_NOT_FOUND'

            utils.pushDataInMap(separatedMovies, moviesData[0].movies, moviesData[0].name);
            utils.pushDataInMap(separatedMovies, moviesData[1].movies, moviesData[1].name);
            for (const [a, b] of separatedMovies) {
                moviesCheck.push(
                    new Promise((resolve, reject) => {
                        utils.fetchDataFromUrl(a).then((element) => {

                            if (element.error) {
                                includingMovies.push(element.url)
                            }
                            else {
                                let hasActor = (b == req.value.actor1) ? req.value.actor2 : req.value.actor1
                                if (element.cast.find((e) => e.name == hasActor)) {
                                    separatedMovies.delete(a);
                                }
                            }
                            resolve()
                        })
                            .catch((err) => {
                                reject(err);
                            })
                    })
                )
            }
            await Promise.all(moviesCheck);
            return utils.sendResponse(res, 200, 'Fetched Successfully', {
                degreeOfSeparation: separatedMovies.size,
                needToIncludeMovies: includingMovies
            });
        } catch (error) {
            let code = 400, message;
            switch (error) {
                case 'ACTOR_MOVIES_URL_ISSUE':
                    message = 'URL ISSUE';
                    break;
                case 'MOVIES_NOT_FOUND':
                    message = 'Requested data not found';
                    break;
                case 'INVALID_URL':
                    message = 'URL ISSUE';
                    break;
                default:
                    code = 500;
                    message = 'Something Went Wrong';
                    break;
            }
            return utils.sendResponse(res, code, message, {});
        }

    }
}