import { makeStyles } from "@material-ui/core/styles";
import AccessTimeIcon from "@material-ui/icons/AccessTime";
import React from "react";

const useStyles = makeStyles(theme => ({
  root: {
    display: "flex",
    "align-items": "center",
    textTransform: "uppercase",
  },
}));

const formatTimestamp = timestamp => {
  if (!timestamp) {
    return "n/a";
  }

  const lastTimestampDate = new Date(timestamp);
  const [weekday, month, day, year, time] = lastTimestampDate
    .toString()
    .split(" ");
  const millis = lastTimestampDate.getMilliseconds();

  return `${weekday} ${month} ${day} ${year} @ ${time}.${millis}`;
};

export default ({ timestamp, className, style }) => {
  const classes = useStyles();

  return (
    <>
      <div
        title="Last Update"
        className={`${classes.root} ${className}`}
        style={style}
      >
        <AccessTimeIcon style={{ marginRight: "4px" }} />
        {formatTimestamp(timestamp)}
      </div>
    </>
  );
};
