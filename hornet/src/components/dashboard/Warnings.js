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
          </TableRow>
        </TableHead>
        <TableBody>
          {items.map(({ timestamp, value: { level, error="-", msg } }, index) => {
            // let TR = trMap[level];
            const date = Date.parse(timestamp)
            const formattedDate = format(date, 'ddd MMM dd YYYY @ HH:mm:ss.SSS').toUpperCase()
            return (
            <TableRow style={{color: `${styleMap[level]}`}} key={`${timestamp}-${index}`} >
              <TableCell style={{color: `${styleMap[level]}`}}>{formattedDate}</TableCell>
              <TableCell style={{color: `${styleMap[level]}`}}>{error}</TableCell>
              <TableCell style={{color: `${styleMap[level]}`}}>{msg}</TableCell>
            </TableRow>
            )
          })}
        </TableBody>
      </Table>
    </>
  );
}
