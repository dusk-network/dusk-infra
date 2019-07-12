/* eslint-disable no-script-url */

import Typography from "@material-ui/core/Typography";
import { makeStyles } from "@material-ui/styles";
import React, { useRef, useState, useEffect } from "react";
import Peak from "./Peak";
import Title from "./Title";

const MAIN_COLOR = "rgba(0, 0, 0, 0.54)";

function useInterval(callback, delay) {
  const savedCallback = useRef();

  // Remember the latest callback.
  useEffect(() => {
    savedCallback.current = callback;
  }, [callback]);

  // Set up the interval.
  useEffect(() => {
    function tick() {
      savedCallback.current();
    }
    if (delay !== null) {
      let id = setInterval(tick, delay);
      return () => clearInterval(id);
    }
  }, [delay]);
}

const useStyles = makeStyles(theme => ({
  blockHash: {
    flex: 1,
  },
  blockHeight: {
    padding: theme.spacing(2),
  },
}));

const fromNow = (now, timestamp) =>
  Math.max(0, (now - new Date(timestamp)) / 1e3).toFixed(1) + "s ago";

export default React.memo(({ timestamp }) => {
  const classes = useStyles();
  let [now, setNow] = useState(Date.now());

  useInterval(() => {
    setNow(Date.now());
  }, 100);

  return (
    <>
      <Title>Last Update</Title>
      <Typography component="p" variant="h4" className={classes.blockHeight}>
        <time dateTime={timestamp}>
          {timestamp ? fromNow(now, timestamp) : "n/a"}
        </time>
      </Typography>
      {/* <LastUpdate
      timestamp={timestamp}
      style={{ color: MAIN_COLOR }}
      /> */}
      <Typography color="textSecondary" noWrap className={classes.blockHash}>
        {timestamp || "n/a"}
      </Typography>
    </>
  );
});
