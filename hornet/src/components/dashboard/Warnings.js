/* eslint-disable no-script-url */
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableHead from "@material-ui/core/TableHead";
import TableRow from "@material-ui/core/TableRow";
import { makeStyles } from "@material-ui/styles";
import { format } from 'date-fns';
import React from "react";
import Title from "./Title";

const styleMap = {
    warning: "#dda458",
    error: "#d70206",
    fatal: "#6b0392",
    panic: "#0544d3"
}

const useStyles = makeStyles(theme => ({}));

export default function Nodes({ items = [] }) {
  const classes = useStyles();

  return (
    <>
      <Title>Last Errors</Title>
      <Table size="small">
        <TableHead>
          <TableRow>
            <TableCell>Timestamp</TableCell>
            <TableCell>Error</TableCell>
            <TableCell>Message</TableCell>
            <TableCell>Process</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {items.map(({ timestamp, value: { process: comp, prefix, level, error="-", msg } }, index) => {
            // let TR = trMap[level];
            const origin = comp ? comp : (prefix ? prefix : "")
            const color = styleMap[level]
            const date = Date.parse(timestamp)
            const formattedDate = format(date, 'ddd MMM dd YYYY @ HH:mm:ss.SSS').toUpperCase()
            return (
            <TableRow style={{color: color}} key={`${timestamp}-${index}`} >
              <TableCell style={{color: color}}>{formattedDate}</TableCell>
              <TableCell style={{color: color}}>{error}</TableCell>
              <TableCell style={{color: color}}>{msg}</TableCell>
              <TableCell style={{color: color}}>{origin}</TableCell>
            </TableRow>
            )
          })}
        </TableBody>
      </Table>
    </>
  );
}
