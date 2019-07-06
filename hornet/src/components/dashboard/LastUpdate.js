import React from "react";
import { ResponsiveContainer } from "recharts";

import AccessTimeIcon from "@material-ui/icons/AccessTime";
import { makeStyles } from "@material-ui/core/styles";

const useStyles = makeStyles(theme => ({
  root: {
    display: "flex",
    "align-items": "center",
    "text-align": "center",
    "border-radius": "4px",
    padding: "4px",
    // color: "rgba(5, 68, 211)",
    // background: "rgba(5, 68, 211, .1)",
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

export default ({ timestamp, className }) => {
  const classes = useStyles();

  return (
    <>
      <div title="Last Update" className={`${classes.root} ${className}`}>
        <AccessTimeIcon style={{ marginRight: "4px" }} />
        {formatTimestamp(timestamp)}
      </div>
    </>
  );
};
