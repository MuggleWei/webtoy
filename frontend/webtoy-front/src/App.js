import { Grid } from "@mui/material";
import React, { useEffect, useState } from "react";
import { Outlet, useLocation, useNavigate } from "react-router-dom";
import HeadMenu from "./components/HeadMenu";
import Navbar from "./components/Navbar";
import UserAnchor from "./components/UserAnchor";
import { functionItems } from "./utils/functionItems";
import requests from "./utils/requests";
import url from "./utils/url";

function App() {
    const [navbarVisable, setNavbarVisable] = useState(false);
    const [useAnchorEl, setUserAnchorEl] = useState(null);
    const [auth, setAuth] = useState(false);

    let navigate = useNavigate();
    let location = useLocation();

    useEffect(() => {
        requests
            .post(url.api.check, {})
            .then(rsp => {
                if (rsp.code) {
                    setAuth(false);
                } else {
                    setAuth(true);
                }
            })
    })

    const handleLogin = () => {
        navigate(url.login, {
            replace: true,
            state: { from: location },
        })
    }

    const handleProfile = () => {
        navigate(url.profile, {
            state: { from: location },
        })
    }

    const handleSettings = () => {
        // TODO:
        console.info("handle setttings");
    }

    const handleLogout = () => {
        sessionStorage.removeItem("uid");
        sessionStorage.removeItem("session");
        sessionStorage.removeItem("token");
        localStorage.removeItem("uid");
        localStorage.removeItem("session");
        localStorage.removeItem("token");
        setAuth(false);
        navigate(url.home, {replace: true});
    }

    return (
        <>
            <Grid container>
                <Navbar
                    functionItems={functionItems}
                    visable={navbarVisable}
                    setVisable={setNavbarVisable}
                />
                <HeadMenu
                    setNavbarVisable={setNavbarVisable}
                    setUserAnchorEl={setUserAnchorEl}
                    auth={auth}
                    handleLogin={handleLogin}
                />
                <UserAnchor
                    anchorEl={useAnchorEl}
                    setAnchorEl={setUserAnchorEl}
                    handleProfile={handleProfile}
                    handleSettings={handleSettings}
                    handleLogout={handleLogout}
                />
            </Grid>
            <Grid container>
                <Outlet />
            </Grid>
        </>
    );
}

export default App;