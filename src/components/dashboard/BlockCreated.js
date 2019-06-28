/* eslint-disable no-script-url */

import React from "react";
import Link from "@material-ui/core/Link";
import { makeStyles } from "@material-ui/styles";
import Typography from "@material-ui/core/Typography";
import Title from "./Title";

const useStyles = makeStyles(
  theme => (
    console.log(theme),
    {
      blockHash: {
        flex: 1
      },
      blockHeight: {
        padding: theme.spacing(2)
      }
    }
  )
);

const fromNow = timestamp =>
  ((Date.now() - new Date(timestamp)) / 1000).toFixed(2) + "s ago";

export default React.memo(({ timestamp }) => {
  const classes = useStyles();
  return (
    <>
      <Title>Last Update</Title>
      <Typography component="p" variant="h4" className={classes.blockHeight}>
        <time dateTime={timestamp}>
          {timestamp ? fromNow(timestamp) : "n/a"}
        </time>
      </Typography>
      <Typography color="textSecondary" noWrap className={classes.blockHash}>
        <time dateTime={timestamp}>{timestamp || "n/a"}</time>
      </Typography>
    </>
  );
});
