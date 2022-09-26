import React from "react";
import { useLocation, useNavigate } from "react-router-dom";
import ErrorCode from "../utils/err-code";
import requests from "../utils/requests";
import url from "../utils/url";

class Profile extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            status: "wait"
        };
    }

    componentDidMount() {
        requests
            .post(url.api.profile)
            .then(rsp => {
                if (!rsp.code || rsp.code === ErrorCode.OK) {
                    this.setState({ status: "success" });
                } else if (rsp.code === ErrorCode.ERROR_AUTH) {
                    console.warn("failed Authentication:", rsp.err_msg)
                    this.props.navigate(url.login, {
                        replace: true,
                        state: { from: this.props.location }
                    });
                } else {
                    console.error(rsp);
                }
            })
    }

    render() {
        if (this.state.status === "wait") {
            return (<></>)
        } else {
            return (
                <>
                    <div>Profile页面</div>
                </>
            )
        }
    }
}


function ProfileHook(props) {
    let navigate = useNavigate();
    let location = useLocation();
    return (
        <Profile {...props} navigate={navigate} location={location} />
    );
}

export default ProfileHook;