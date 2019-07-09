/* eslint-disable no-script-url */

import Typography from "@material-ui/core/Typography";
import { makeStyles } from "@material-ui/styles";
import React from "react";
import LastUpdate from './LastUpdate';
import Title from "./Title";

const MAIN_COLOR = "rgba(0, 0, 0, 0.54)";

const useStyles = makeStyles(theme => ({
  blockHash: {
    flex: 1,
  },
  blockHeight: {
    padding: theme.spacing(2),
  },
}));

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
      <LastUpdate
      timestamp={timestamp}
      style={{ color: MAIN_COLOR }}
      />
      {/* <Typography color="textSecondary" noWrap className={classes.blockHash}>
      </Typography> */}
    </>
  );
});
