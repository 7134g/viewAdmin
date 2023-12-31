import axios from "axios";

const request = axios.create({
    timeout: 5000
});

request.interceptors.request.use(
    config => {
        config.headers["Content-Type"] = "application/json";
        return config;
    },
    error => {
        return Promise.reject(error);
    }
);

request.interceptors.response.use(
    response => {
        let res = response.data;
        if (response.config.responseType === "blob") {
            return res;
        }

        if (typeof res == "string") {
            res = res ? JSON.parse(res) : res;
        }
        return res;
    },
    error => {
        console.log("err" + error);
        return Promise.reject(error);
    }
);

export default request;
