import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import App from "./App";
import Home from "./pages/Home";
import Login from "./pages/Login";
import Profile from "./pages/Profile";
import url from "./utils/url";
import NotFound from "./pages/NotFound"

const root = ReactDOM.createRoot(document.getElementById("root"));
root.render(
  <BrowserRouter>
    <Routes>
      <Route path={url.login} element={<Login />} />
      <Route path={url.home} element={<App />}>
        <Route index element={<Home />} />
        <Route path={url.profile} exact element={<Profile />} />
      </Route>
      <Route path="*" element={<NotFound />} />
    </Routes>
  </BrowserRouter>
);