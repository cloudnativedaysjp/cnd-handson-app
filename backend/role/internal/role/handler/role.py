from role import role_pb2_grpc, role_pb2
import grpc


class RoleHandler(role_pb2_grpc.RoleServiceServicer):
    """RoleHandlerの実装"""

    def __init__(self, role_service):
        """RoleHandlerの初期化
        Args:
            role_service: RoleServiceのインスタンス
        """
        self.role_service = role_service

    def CreateRole(self, request, context):
        """Roleの生成
        Args:
        request: CreateRoleRequest
            Roleの生成に必要な情報を含むリクエスト
        context: grpc.ServicerContext
        """
        try:
            # RoleService から生成された RoleModel を取得
            role_model = self.role_service.create(request.name, request.description)

            # RoleResponse を構築
            response = role_pb2.RoleResponse(
                role=role_pb2.Role(
                    id=str(role_model.id),
                    name=role_model.name,
                    description=role_model.description,
                )
            )
            return response
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Error creating role: {str(e)}")
            return role_pb2.RoleResponse()

    def GetRole(self, request, context):
        """Roleの取得"""
        try:
            # RoleService から RoleModel を取得
            role_model = self.role_service.get(role_id=request.id)
            if role_model is None:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                context.set_details("Role not found!")
                return role_pb2.RoleResponse()

            # RoleResponse を構築
            response = role_pb2.RoleResponse(
                role=role_pb2.Role(
                    id=str(role_model.id),
                    name=role_model.name,
                    description=role_model.description,
                )
            )
            return response
        except Exception as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"Role not found: {str(e)}")
            return role_pb2.RoleResponse()

    def UpdateRole(self, request, context):
        """Roleの更新"""
        try:

            # ここでRoleを更新する処理を実装する
            role_model = self.role_service.update(
                request.id, request.name, request.description
            )
            if role_model is None:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                context.set_details("Role not found!")
                return role_pb2.RoleResponse()

            # RoleResponse を構築
            response = role_pb2.RoleResponse(
                role=role_pb2.Role(
                    id=str(role_model.id),
                    name=role_model.name,
                    description=role_model.description,
                )
            )
            return response
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Error updating role: {str(e)}")
            return role_pb2.RoleResponse()

    def DeleteRole(self, request, context):
        """Roleの削除"""
        # ここでRoleを削除する処理を実装する
        success = self.role_service.delete(request.id)
        if not success:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("Role not found!")
            return role_pb2.DeleteRoleResponse(success=False)
        response = role_pb2.DeleteRoleResponse(success=True)
        return response

    def ListRoles(self, request, context):
        """Roleの一覧取得"""
        # ここでRoleの一覧を取得する処理を実装する
        roles = self.role_service.list(request, request.page, request.page_size)
        response = role_pb2.ListRolesResponse(roles=roles)
        return response
