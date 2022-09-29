import Cookies from 'js-cookie';
import React, { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import Sign from "../components/Sign"
import ErrorCode from '../utils/errCode';
import requests from '../utils/requests';
import url from '../utils/url';

function Login() {
    const [servRsped, setServRsped] = useState(false);
    const [authError, setAuthError] = useState("");
    const [captchaError, setCaptchaError] = useState("");
    const [captchaLoadTime, setCaptchaLoadTime] = useState(Date.now());
    const [rememberMe, setRememberMe] = useState(false);

    let navigate = useNavigate();
    let location = useLocation();

    let from = location.state?.from?.pathname || url.home;

    useEffect(() => {
        requests
            .post(url.api.check, {})
            .then(rsp => {
                if (rsp.code) {
                    setServRsped(true);
                } else {
                    navigate(from, { replace: true });
                }
            })
    })

    const handleSubmit = (e) => {
        e.preventDefault();
        const data = new FormData(e.currentTarget);

        setCaptchaError("");
        setAuthError("");

        // login
        requests.post(url.api.login, {
            name: data.get("user"),
            passwd: data.get("password"),
            captcha_session: Cookies.get("captcha_session"),
            captcha_value: data.get("captcha"),
        }).then(rsp => {
            if (!rsp.code) {
                // save session & token
                sessionStorage.setItem("uid", rsp.data.uid);
                sessionStorage.setItem("session", rsp.data.session);
                sessionStorage.setItem("token", rsp.data.token);
                if (rememberMe) {
                    localStorage.setItem("uid", rsp.data.uid);
                    localStorage.setItem("session", rsp.data.session);
                    localStorage.setItem("token", rsp.data.token);
                }

                navigate(from, { replace: true });
            } else {
                console.error(rsp.msg || "login failed");
                if (rsp.code === ErrorCode.ERROR_CAPTCHA) {
                    setCaptchaError("验证码错误");
                } else if (rsp.code === ErrorCode.ERROR_AUTH) {
                    setAuthError("用户名或密码错误");
                } else {
                    setAuthError("登录异常");
                }
                setCaptchaLoadTime(Date.now());
            }
        })
    }

    return (servRsped
        ? <Sign
            urlCaptcha={url.api.captcha}
            authError={authError}
            captchaError={captchaError}
            rememberMe={rememberMe}
            setRememberMe={setRememberMe}
            captchaLoadTime={captchaLoadTime}
            setCaptchaLoadTime={setCaptchaLoadTime}
            handleSubmit={handleSubmit}
        />
        : <></>
    );
}

export default Login;