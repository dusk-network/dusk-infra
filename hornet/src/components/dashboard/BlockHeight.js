/* eslint-disable no-script-url */

import React from "react";
import Link from "@material-ui/core/Link";
import { makeStyles } from "@material-ui/styles";
import Typography from "@material-ui/core/Typography";
import Title from "./Title";

const useStyles = makeStyles(theme => ({
  blockHash: {
    flex: 1,
    textTransform: "uppercase"
  },
  blockHeight: {
    padding: theme.spacing(2)
  }
}));

export default React.memo(({ height, hash }) => {
  const classes = useStyles();
  return (
    <React.Fragment>
      <Title>Current Block Height</Title>
      <Typography component="p" variant="h4" className={classes.blockHeight}>
        # {height ? height.toLocaleString("en-US") : "n/a"}
      </Typography>
      <Typography color="textSecondary" noWrap className={classes.blockHash}>
        {hash || "n/a"}
      </Typography>
    </React.Fragment>
  );
});
