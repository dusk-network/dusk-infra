import React from "react";
import { ResponsiveContainer } from "recharts";
import { makeStyles } from "@material-ui/styles";
import tenor from "../../tenor.gif";
import ChartistGraph from "react-chartist";

import Title from "./Title";

const options = {};

const type = "Pie";

const useStyles = makeStyles(theme => ({
  tenor: {
    width: "100%",
    height: "100%",
    backgroundImage: `url(${tenor})`,
    backgroundSize: "contain",
    backgroundRepeat: "no-repeat",
    backgroundPosition: "center",
  },
}));
export default ({ data }) => {
  const classes = useStyles();

  return (
    <>
      <Title>Disk Usage (%)</Title>
      <ResponsiveContainer>
        {data ? (
          <ChartistGraph data={data} type={type} options={options} />
        ) : (
          <div className={classes.tenor} />
        )}
      </ResponsiveContainer>
    </>
  );
};
