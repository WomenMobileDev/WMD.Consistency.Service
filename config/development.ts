
module.exports = {
    port: 3000,
    enableFileLogs: true,
    enableConsoleLogs: true,
    googleOauth: {
        clientId: "",
        clientSecret: "",
        callbackURL: "",
    },mysqlRemote: {
        host: "",
        port: "",
        username: "",
        password: "",
        db: "",
    },
    mysqlLocal: {
        host: "",
        username: "",
        password: "",
        db: "",
    },
    services: {
        mtApi: {
            baseUrl: "",
        },

        mtUi: {
            baseUrl: ""
        }
    }
};