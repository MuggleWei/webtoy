import { Container } from "@mui/material";
import React, { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import requests from "../utils/requests";
import url from "../utils/url";

function Profile() {
    let navigate = useNavigate();
    let location = useLocation();

    const [userProfileInfo, setUserProfileInfo] = useState(null);

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
                    console.info(rsp);
                    setUserProfileInfo(rsp.data);
                }
            })
    }, [userProfileInfo, navigate, location])

    return (userProfileInfo
        ?
        <Container>
            <div>
                {userProfileInfo.show_name}
            </div>
            <div>
                <h3>User ID: {userProfileInfo.uid}</h3>
                <h3>Email: {userProfileInfo.email}</h3>
                <h3>Phone: {userProfileInfo.phone}</h3>
            </div>
        </Container>
        :
        <></>
    );
}

export default Profile;