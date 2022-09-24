import { Grid } from '@mui/material';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Checkbox from '@mui/material/Checkbox';
import Container from '@mui/material/Container';
import CssBaseline from '@mui/material/CssBaseline';
import FormControlLabel from '@mui/material/FormControlLabel';
import Link from '@mui/material/Link';
import { createTheme, ThemeProvider } from '@mui/material/styles';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';
import Cookies from 'js-cookie';
import React from 'react';
import { useNavigate } from 'react-router-dom';

import ErrorCode from '../utils/err-code';
import requests from "../utils/requests";
import url from "../utils/url";

function Copyright(props) {
    return (
        <Typography variant="body2" color="text.secondary" align="center" {...props}>
            {'Copyright © '}
            <Link color="inherit" href="https://github.com/MuggleWei/webtoy">
                muggle wei
            </Link>{' '}
            {new Date().getFullYear()}
            {'.'}
        </Typography>
    );
}

const theme = createTheme();

class Login extends React.Component {
    constructor(props) {
        super(props);

        this.urlCaptcha = url.api.captchaLoad;

        this.state = {
            captchaLoadTime: Date.now(),
            rememberMeChecked: false,
        }

        this.handleSubmit = this.handleSubmit.bind(this);
        this.handleRememberMeChange = this.handleRememberMeChange.bind(this);
        this.handleCaptchaClick = this.handleCaptchaClick.bind(this);
    }

    handleSubmit(e) {
        e.preventDefault();
        const data = new FormData(e.currentTarget);

        // login
        requests.post(url.api.login, {
            user: data.get("user"),
            password: data.get("password"),
            captcha_session: Cookies.get("captcha_session"),
            captcha_value: data.get("captcha"),
        }).then(rsp => {
            if (!rsp.code || rsp.code === ErrorCode.OK) {
                // save session & token
                sessionStorage.setItem("uid", rsp.data.uid);
                sessionStorage.setItem("session", rsp.data.session);
                sessionStorage.setItem("token", rsp.data.token);
                if (this.state.rememberMeChecked) {
                    localStorage.setItem("uid", rsp.data.uid);
                    localStorage.setItem("session", rsp.data.session);
                    localStorage.setItem("token", rsp.data.token);
                }

                // navigate to home
                this.props.navigate(url.home, { replace: true });
            } else {
                console.error(rsp.msg || "login failed");
            }
        })
    };

    handleRememberMeChange(e) {
        this.setState({ rememberMeChecked: e.target.checked });
    }

    handleCaptchaClick() {
        this.setState({
            captchaLoadTime: Date.now(),
        })
    }

    render() {
        return (
            <ThemeProvider theme={theme}>
                <Container component="main" maxWidth="xs">
                    <CssBaseline />
                    <Grid
                        sx={{
                            marginTop: 8,
                            display: 'flex',
                            flexDirection: 'column',
                            alignItems: 'center',
                        }}
                    >
                        <Typography component="h1" variant="h5">
                            Sign in
                        </Typography>
                        <Grid component="form" onSubmit={this.handleSubmit} noValidate sx={{ mt: 1 }}>
                            <TextField
                                margin="normal"
                                required
                                fullWidth
                                id="user"
                                name="user"
                                label="User"
                                autoComplete="user"
                                autoFocus
                            />
                            <TextField
                                margin="normal"
                                required
                                fullWidth
                                id="password"
                                name="password"
                                label="Password"
                                type="password"
                                autoComplete="current-password"
                            />
                            <Grid container>
                                <Grid container sx={{ width: 5 / 10 }}>
                                    <TextField
                                        margin="normal"
                                        required
                                        id="captcha"
                                        name="captcha"
                                        label="Captcha"

                                    />
                                </Grid>
                                <Grid sx={{ width: 1 / 10 }}>
                                </Grid>
                                <Box onClick={this.handleCaptchaClick}
                                    sx={{ width: 4 / 10 }}
                                    component="img"
                                    src={`${this.urlCaptcha}?${this.state.captchaLoadTime}`}
                                    alt="captcha"
                                />
                            </Grid>
                            <FormControlLabel
                                control={<Checkbox value="remember" color="primary" checked={this.state.rememberMeChecked} onChange={this.handleRememberMeChange} />}
                                label="Remember me"
                                name="rememberme"
                            />
                            <Button
                                type="submit"
                                fullWidth
                                variant="contained"
                                sx={{ mt: 3, mb: 2 }}
                            >
                                Sign In
                            </Button>
                        </Grid>
                    </Grid>
                    <Copyright sx={{ mt: 8, mb: 4 }} />
                </Container>
            </ThemeProvider>
        );
    }
}

function LoginWithNavigation(props) {
    let navigate = useNavigate();
    return (
        <Login {...props} navigate={navigate} />
    );
}

export default LoginWithNavigation;