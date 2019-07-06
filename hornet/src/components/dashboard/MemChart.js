import { makeStyles } from "@material-ui/styles";
import React from "react";
import ChartistGraph from "react-chartist";
import { ResponsiveContainer } from "recharts";
import * as chartUtils from "../../chart-utils";
import LastUpdate from "./LastUpdate";
import Title from "./Title";



const options = {
  fullWidth: true,
  showArea: true,
  chartPadding: {
    right: 40,
  },
  high: 100,
  low: 0,
  showPoint: true,
  lineSmooth: true,
  axisX: {
    labelInterpolationFnc: chartUtils.skipLabels,
  },
};

const type = "Line";
const useStyles = makeStyles(theme => ({
  lastUpdate: {
    color: "#D70206",
  },
}));

export default ({ data }) => {
  const classes = useStyles();

  return (
    <>
      <Title>Memory Usage (%)</Title>
      <ResponsiveContainer>
        <ChartistGraph
          data={data}
          type={type}
          options={options}
          listener={chartUtils.listener("mem-timestamp")}
        />
      </ResponsiveContainer>
      <LastUpdate
        timestamp={data.labels[data.labels.length - 1]}
        className={classes.lastUpdate}
      />
    </>
  );
};
