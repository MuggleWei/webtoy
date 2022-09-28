import AccountCircle from '@mui/icons-material/AccountCircle';
import MenuIcon from '@mui/icons-material/Menu';
import { Box } from "@mui/material";
import AppBar from '@mui/material/AppBar';
import Button from '@mui/material/Button';
import IconButton from '@mui/material/IconButton';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import React from 'react';

function HeadMenu({
    setNavbarVisable,
    setUserAnchorEl,
    auth = false,
    handleLogin
}) {
    return (
        <Box sx={{ flexGrow: 1 }}>
            <AppBar position="static">
                <Toolbar>
                    {/* 菜单栏左侧抽屉按钮 */}
                    <IconButton
                        size="large"
                        edge="start"
                        color="inherit"
                        aria-label="menu"
                        sx={{ mr: 2 }}
                        onClick={() => setNavbarVisable(true)}
                    >
                        <MenuIcon />
                    </IconButton>
                    {/* 菜单栏左侧文字说明 */}
                    <Typography
                        variant="h6"
                        noWrap
                        component="div"
                        sx={{ flexGrow: 1, display: { xs: 'none', sm: 'block' } }}
                    >
                        Menu
                    </Typography>
                    {/* 用户信息或登录按钮 */}
                    {auth
                        ? <>
                            <IconButton
                                size="large"
                                aria-label="account of current user"
                                aria-controls="menu-appbar"
                                aria-haspopup="true"
                                onClick={(e) => setUserAnchorEl(e.currentTarget)}
                                color="inherit"
                            >
                                <AccountCircle />
                            </IconButton>
                        </>
                        : <Button
                            color="inherit"
                            onClick={handleLogin}
                        >
                            Login
                        </Button>}
                </Toolbar>
            </AppBar>
        </Box>
    );
}

export default HeadMenu;