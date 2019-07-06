/* eslint-disable no-script-url */

import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableHead from "@material-ui/core/TableHead";
import TableRow from "@material-ui/core/TableRow";
import { makeStyles } from "@material-ui/styles";
import { format } from "date-fns";
import React from "react";

const useStyles = makeStyles(theme => ({}));

export default function Nodes({ items = [] }) {
  const classes = useStyles();
  return (
    <>
      <Table size="small">
        <TableHead>
          <TableRow>
            <TableCell>Timestamp</TableCell>
            <TableCell>Type</TableCell>
            <TableCell>Message</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {items.map(({ timestamp, value }, index) => {
            const date = Date.parse(timestamp)
            const formattedDate = format(date, 'ddd MMM dd YYYY @ HH:mm:ss.SSS').toUpperCase()
            return (
            <TableRow key={`${timestamp}-${index}`}>
              <TableCell>{formattedDate}</TableCell>
              <TableCell>{value}</TableCell>
            </TableRow>
            )
          })}
        </TableBody>
      </Table>
    </>
  );
}
