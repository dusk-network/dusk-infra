import React from "react";
import PropTypes from "prop-types";
import Typography from "@material-ui/core/Typography";

const Title = props => (
  <Typography component="p" variant="body2" color="primary" gutterBottom>
    {props.children}
  </Typography>
);

Title.propTypes = {
  children: PropTypes.node,
};

export default Title;
