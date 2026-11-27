"""Modbus TCP simulator for local GoPLC development.

Serves holding registers on 0.0.0.0:5020 and answers any unit ID.
Seeded holding registers (0-based protocol addresses):

    0..2    100, 200, 300                three plain uint16 values
    10..11  0x41AC, 0x0000               float32 21.5, ABCD order
    12..13  0xFFFF, 0xFC18               int32 -1000, ABCD order
    14      0xFFFB                       int16 -5

Coils 0..7 alternate on/off. Input registers 0..2 hold 1, 2, 3.
"""

import asyncio
import logging

from pymodbus.datastore import (
    ModbusDeviceContext,
    ModbusSequentialDataBlock,
    ModbusServerContext,
)
from pymodbus.server import StartAsyncTcpServer

PORT = 5020

HOLDING = [0] * 32
HOLDING[0:3] = [100, 200, 300]
HOLDING[10:12] = [0x41AC, 0x0000]
HOLDING[12:14] = [0xFFFF, 0xFC18]
HOLDING[14] = 0xFFFB


def build_context() -> ModbusServerContext:
    # Device contexts add 1 to the protocol address, so blocks start at 1
    # for protocol address 0 to map to the first value.
    device = ModbusDeviceContext(
        hr=ModbusSequentialDataBlock(1, HOLDING),
        ir=ModbusSequentialDataBlock(1, [1, 2, 3] + [0] * 29),
        co=ModbusSequentialDataBlock(1, [True, False] * 16),
        di=ModbusSequentialDataBlock(1, [False] * 32),
    )
    return ModbusServerContext(devices=device, single=True)


async def main() -> None:
    logging.basicConfig(level=logging.INFO)
    logging.info("Modbus TCP simulator listening on 0.0.0.0:%d", PORT)
    await StartAsyncTcpServer(context=build_context(), address=("0.0.0.0", PORT))


if __name__ == "__main__":
    asyncio.run(main())
