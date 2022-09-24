import React from "react";
import { useNavigate } from "react-router-dom";
import ErrorCode from "../utils/err-code";
import requests from "../utils/requests";
import url from "../utils/url";

class ProfileWithNavigate extends React.Component {
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
                    this.props.navigate(url.login, { replace: true });
                } else {
                    console.error(rsp);
                }
            })
    }

    render() {
        if (this.state.status === "wait") {
            return (<div>页面跳转中</div>)
        } else {
            return (
                <>
                    <div>Profile页面</div>
                </>
            )
        }
    }
}


function Profile(props) {
    let navigate = useNavigate();
    return (
        <ProfileWithNavigate {...props} navigate={navigate} />
    );
}

export default Profile;