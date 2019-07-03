import React from "react";
import PropTypes from "prop-types";
import Typography from "@material-ui/core/Typography";
import { makeStyles } from "@material-ui/styles";

const useStyles = makeStyles(theme => ({
  root: {
    textAlign: "center",
    transform: "translateY(-100%)",
  },
}));

export default function LastUpdate(props) {
  const classes = useStyles();

  return (
    <Typography
      className={classes.root}
      component="p"
      variant="body2"
      color="primary"
    >
      {props.children}
    </Typography>
  );
}

Description.propTypes = {
  children: PropTypes.node,
};
