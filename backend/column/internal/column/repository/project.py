from uuid import UUID

import grpc

from pkg.userid import userid
from project import project_pb2, project_pb2_grpc


class ProjectRepository:
    """project に、呼び出し元がプロジェクトの所有者かを問い合わせる"""

    def __init__(self, channel: grpc.Channel):
        self.stub = project_pb2_grpc.ProjectServiceStub(channel)

    def check_access(self, project_id: UUID) -> bool:
        try:
            self.stub.CheckProjectAccess(
                project_pb2.CheckProjectAccessRequest(project_id=str(project_id)),
                metadata=userid.forward_metadata(),
                timeout=5,
            )
        except grpc.RpcError as e:
            if e.code() == grpc.StatusCode.NOT_FOUND:
                return False
            raise
        return True
