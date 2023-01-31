import { Badge, Box, Center } from "@mantine/core";
import { XCircleIcon } from "@primer/octicons-react";
import { forwardRef } from "react";

export const ReportTag = forwardRef<HTMLDivElement, { name: string; color: string; onRemove?: () => void }>(
  ({ name, color, onRemove, ...others }, ref) => (
    <div ref={ref} {...others}>
      <Badge color={color}>
        <Box
          sx={{
            display: "flex",
          }}
        >
          <p>{name}</p>
          {onRemove && (
            <Center
              ml={5}
              sx={{
                cursor: "pointer",
              }}
              onMouseDown={onRemove}
            >
              <XCircleIcon size={14} />
            </Center>
          )}
        </Box>
      </Badge>
    </div>
  )
);
