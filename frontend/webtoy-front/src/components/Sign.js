import { Button, Checkbox, Container, CssBaseline, FormControlLabel, Grid, TextField, Typography } from "@mui/material";
import Captcha from "./Captcha";
import Copyright from "./Copyright";

function Sign({ urlCaptcha, authError, captchaError, rememberMe, setRememberMe, captchaLoadTime, setCaptchaLoadTime, handleSubmit }) {
    return (
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
                <Grid component="form" onSubmit={handleSubmit} noValidate sx={{ mt: 1 }}>
                    <TextField
                        margin="normal"
                        required
                        fullWidth
                        id="user"
                        name="user"
                        label="User"
                        autoComplete="user"
                        autoFocus
                        error={authError ? true : false}
                        helperText={authError}
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
                        error={authError ? true : false}
                        helperText={authError}
                    />
                    <Captcha
                        id="captcha"
                        url={urlCaptcha}
                        captchaError={captchaError} 
                        loadTime={captchaLoadTime}
                        setLoadTime={setCaptchaLoadTime}
                        />
                    <FormControlLabel
                        control={
                            <Checkbox
                                value="remember"
                                color="primary"
                                checked={rememberMe}
                                onChange={(e) => setRememberMe(e.target.checked)} />
                        }
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
    );
}

export default Sign;