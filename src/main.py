import asyncio
import logging
import sys

import grpc

sys.path.insert(0, "./pb")

from concurrent import futures

from rembg import remove

from pb.cleaner_pb2 import CleanRequest, CleanResponse
from pb.cleaner_pb2_grpc import (CleanerServiceServicer,
                                 add_CleanerServiceServicer_to_server)

_cleanup_coroutines = []


class Cleaner(CleanerServiceServicer):

    async def Clean(
        self,
        request: CleanRequest,
        context: grpc.aio.ServicerContext,
    ) -> CleanResponse:
        output = remove(request.data)
        return CleanResponse(data=output)


async def serve() -> None:
    server = grpc.aio.server(
        futures.ThreadPoolExecutor(max_workers=10),
        # options=[
        #     ("grpc.max_receive_message_length", 32 << 20),
        #     ("grpc.max_send_message_length", 32 << 20),
        # ],
    )
    add_CleanerServiceServicer_to_server(Cleaner(), server)
    listen_addr = "[::]:42069"
    server.add_insecure_port(listen_addr)
    await server.start()

    async def server_graceful_shutdown():
        logging.info("Staring graceful shutdown...")
        await server.stop(5)

    _cleanup_coroutines.append(server_graceful_shutdown())
    await server.wait_for_termination()


if __name__ == "main":
    logging.basicConfig(level=logging.INFO)
    loop = asyncio.get_event_loop()
    try:
        loop.run_until_complete(serve())
    finally:
        loop.run_until_complete(*_cleanup_coroutines)
        loop.close()
