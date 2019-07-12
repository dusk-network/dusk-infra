import Typography from "@material-ui/core/Typography";
import { makeStyles } from "@material-ui/core/styles";
import TrendingIcon from "@material-ui/icons/TrendingUp";
import AccessTimeIcon from "@material-ui/icons/AccessTime";
import React, { useState } from "react";
import Tooltip from "@material-ui/core/Tooltip";
import { withStyles } from "@material-ui/core/styles";

const useStyles = makeStyles(theme => ({
  root: {
    display: "flex",
    "align-items": "center",
    "& b": {
      margin: "0 .3em",
      padding: "2px .5em",
      borderRadius: "1em",
      color: "#fff",
      border: "2px solid currentColor",
      cursor: "default",
    },
  },
}));

export const formatTimestamp = timestamp => {
  if (!timestamp) {
    return "n/a";
  }

  const lastTimestampDate = new Date(timestamp);
  const [, , , , time] = lastTimestampDate.toString().split(" ");
  const millis = lastTimestampDate.getMilliseconds();

  return `${time}.${millis}`;
};

const HtmlTooltip = withStyles(theme => ({
  tooltip: {
    backgroundColor: "#f5f5f9",
    color: "rgba(0, 0, 0, 0.87)",
    maxWidth: 220,
    fontSize: theme.typography.pxToRem(12),
    border: "1px solid #dadde9",
  },
}))(Tooltip);

export default ({ timestamp, value, unit, className, style }) => {
  if (!timestamp) {
    return null;
  }
  const numeric = +value;
  const [peak, setPeak] = useState({ numeric, timestamp });
  const classes = useStyles();

  if (numeric > peak.numeric) {
    setPeak({ numeric, timestamp });
  }

  return (
    <>
      <div className={`${classes.root} ${className}`} style={style}>
        <TrendingIcon style={{ marginRight: "4px" }} />
        {"Peak "}
        <HtmlTooltip
          disableFocusListener
          title={
            <div style={{ display: "flex" }}>
              <AccessTimeIcon style={{ marginRight: "4px" }} />

              <Typography color="inherit">
                {formatTimestamp(peak.timestamp)}
              </Typography>
            </div>
          }
        >
          <b style={{ ...style, borderColor: style.color }}>
            {peak.numeric}
            {unit}
          </b>
        </HtmlTooltip>
      </div>
    </>
  );
};
