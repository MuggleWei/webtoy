import MenuIcon from '@mui/icons-material/Menu';
import { Box } from "@mui/material";
import Divider from '@mui/material/Divider';
import Drawer from '@mui/material/Drawer';
import IconButton from '@mui/material/IconButton';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import React from 'react';
import { useNavigate } from "react-router-dom";

function Navbar({functionItems={}, visable=false, setVisable, width=250}) {
    const navigate = useNavigate();

    const handleKeyDown = () => (event) => {
        if (event.key === 'Esc') {
            setVisable(false);
        }
    }

    return (
        <Drawer
            open={visable}
            onKeyDown={handleKeyDown()}
            onClose={() => setVisable(false)}
        >
            <Box
                sx={{ width: width }}
                role="presentation"
            >
                <Toolbar>
                    <IconButton
                        size="large"
                        edge="start"
                        color="inherit"
                        aria-label="navbar"
                        sx={{ mr: 2 }}
                        onClick={() => setVisable(false)}
                    >
                        <MenuIcon />
                    </IconButton>
                    <Typography
                        variant="h6"
                        noWrap
                        component="div"
                        sx={{ flexGrow: 1, display: { xs: 'none', sm: 'block' } }}
                    >
                        Hello World
                    </Typography>
                </Toolbar>
            </Box>
            <Divider />
            <List>
                {functionItems.map((item, index) => (
                    <ListItem
                        button
                        key={item.id}
                        onClick={() => {
                            setVisable(false);
                            navigate(item.route);
                        }}>
                        <ListItemIcon>
                            {item.icon}
                        </ListItemIcon>
                        <ListItemText primary={item.label} />
                    </ListItem>
                ))}
            </List>
        </Drawer>
    );
}

export default Navbar;