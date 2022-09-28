import { Box, Grid, TextField } from "@mui/material";

function Captcha({ id, url, captchaError, loadTime, setLoadTime }) {
    return (
        <Grid container>
            <Grid container sx={{ width: 5 / 10 }}>
                <TextField
                    margin="normal"
                    required
                    id={id}
                    name="captcha"
                    label="Captcha"
                    error={captchaError ? true : false}
                    helperText={captchaError}
                />
            </Grid>
            <Grid sx={{ width: 1 / 10 }}>
            </Grid>
            <Box onClick={() => setLoadTime(Date.now())}
                sx={{ width: 4 / 10 }}
                component="img"
                src={`${url}?${loadTime}`}
                alt="captcha"
            />
        </Grid>
    );
}

export default Captcha;