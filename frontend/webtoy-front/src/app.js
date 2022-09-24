import React from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Home from "./pages/home"
import Login from "./pages/login"
import NotFound from "./pages/not-found";
import Profile from "./pages/profile";
import url from "./utils/url"

class App extends React.Component {
    render() {
        return (
            <BrowserRouter>
                <Routes>
                    <Route path={url.home} element={<Home />} />
                    <Route path={url.login} element={<Login />} />
                    <Route path={url.profile} element={<Profile />} />
                    <Route path="*" element={<NotFound />} />
                </Routes>
            </BrowserRouter>
        );
    }
}

export default App;