import React from "react";
import {
  LineChart,
  CartesianGrid,
  Line,
  XAxis,
  YAxis,
  Label,
  Tooltip,
  Legend,
  ResponsiveContainer
} from "recharts";
import Title from "./Title";

export default ({ data }) => (
  <>
    <Title>Block Time</Title>
    <ResponsiveContainer>
      <LineChart
        data={data}
        margin={{
          top: 16,
          right: 16,
          bottom: 0,
          left: 24
        }}
      >
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="height" />
        <YAxis dataKey="time" />
        <Tooltip />
        <Legend />
        <Line type="monotone" dataKey="time" stroke="#556CD6" dot={false} />
      </LineChart>
    </ResponsiveContainer>
  </>
);
