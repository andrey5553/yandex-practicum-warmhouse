from pydantic import BaseModel
from typing import Optional, Dict, Any, List
from datetime import datetime

class DeviceBase(BaseModel):
    name: str
    type: str
    device_type_id: int
    room_id: Optional[int] = None
    serial_number: str
    configuration: Optional[Dict[str, Any]] = None

class DeviceCreate(DeviceBase):
    pass

class DeviceUpdate(BaseModel):
    name: Optional[str] = None
    configuration: Optional[Dict[str, Any]] = None

class DeviceResponse(DeviceBase):
    id: int
    status: str
    last_seen: Optional[datetime] = None
    created_at: datetime

    class Config:
        from_attributes = True  # Это критически важно!

class DeviceListResponse(BaseModel):
    devices: List[DeviceResponse]

class DeviceCommand(BaseModel):
    command: str
    parameters: Optional[Dict[str, Any]] = None
    priority: str = "normal"

class CommandResult(BaseModel):
    command_id: str
    status: str
    message: Optional[str] = None
    timestamp: datetime

class DeviceState(BaseModel):
    device_id: int
    state: Dict[str, Any]
    last_updated: datetime