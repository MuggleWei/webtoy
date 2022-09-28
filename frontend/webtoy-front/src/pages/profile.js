import { Box } from "@mui/material";
import React, { useEffect } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import requests from "../utils/requests";
import url from "../utils/url";

function Profile() {
    let navigate = useNavigate();
    let location = useLocation();

    useEffect(() => {
        requests
            .post(url.api.profile, {})
            .then(rsp => {
                if (rsp.code) {
                    navigate(url.login, {
                        replace: true,
                        state: { from: location }
                    });
                } else {
                    // TODO: load profile informations
                }
            })
    })

    return (
        <Box>
            Profile Page
        </Box>
    );
}

export default Profile;