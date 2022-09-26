import React from "react";
import { Link, useLocation } from "react-router-dom";
import url from "../utils/url";

class Home extends React.Component {
    render() {
        return (
            <>
                <div>主页</div>
                <Link to={url.login} state={{ from: this.props.location }}>登录</Link>
            </>
        )
    }
}

function HomeHook(props) {
    return (
        <Home {...props} location={useLocation()} />
    )
}

export default HomeHook;