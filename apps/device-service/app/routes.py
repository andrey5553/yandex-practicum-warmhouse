from fastapi import APIRouter, Depends, HTTPException, Query
from sqlalchemy.orm import Session
from typing import List, Optional
from . import schemas, models, database

router = APIRouter()

# In-memory storage for demo (replace with real database in production)
devices_db = {}
command_results = {}

@router.get("/devices", response_model=schemas.DeviceListResponse)
def get_devices(
    room_id: Optional[int] = Query(None),
    status: Optional[str] = Query(None),
    db: Session = Depends(database.get_db)
):
    """Получить список устройств"""
    try:
        query = db.query(models.Device).filter(models.Device.is_active == True)
        
        if room_id:
            query = query.filter(models.Device.room_id == room_id)
        
        if status:
            query = query.filter(models.Device.status == status)
        
        devices = query.all()
        
        # Преобразуем в Pydantic схемы
        device_schemas = [schemas.DeviceResponse.from_orm(device) for device in devices]
        
        return {"devices": device_schemas}
    except Exception as e:
        print(f"Error: {e}")
        raise

@router.post("/devices", response_model=schemas.DeviceResponse, status_code=201)
def create_device(device: schemas.DeviceCreate, db: Session = Depends(database.get_db)):
    """Зарегистрировать новое устройство"""
    # Check if device with serial number already exists
    existing_device = db.query(models.Device).filter(
        models.Device.serial_number == device.serial_number
    ).first()
    
    if existing_device:
        raise HTTPException(status_code=400, detail="Device with this serial number already exists")
    
    db_device = models.Device(
        name=device.name,
        type=device.type,
        device_type_id=device.device_type_id,
        room_id=device.room_id,
        serial_number=device.serial_number,
        configuration=device.configuration or {},
        status="online"
    )
    
    db.add(db_device)
    db.commit()
    db.refresh(db_device)
    
    return db_device

@router.get("/devices/{device_id}", response_model=schemas.DeviceResponse)
def get_device(device_id: int, db: Session = Depends(database.get_db)):
    """Получить информацию об устройстве"""
    device = db.query(models.Device).filter(
        models.Device.id == device_id,
        models.Device.is_active == True
    ).first()
    
    if not device:
        raise HTTPException(status_code=404, detail="Device not found")
    
    return device

@router.put("/devices/{device_id}")
def update_device(device_id: int, device_update: schemas.DeviceUpdate, db: Session = Depends(database.get_db)):
    """Обновить информацию об устройстве"""
    device = db.query(models.Device).filter(
        models.Device.id == device_id,
        models.Device.is_active == True
    ).first()
    
    if not device:
        raise HTTPException(status_code=404, detail="Device not found")
    
    update_data = device_update.dict(exclude_unset=True)
    for field, value in update_data.items():
        setattr(device, field, value)
    
    db.commit()
    
    return {"message": "Device updated successfully"}

@router.post("/devices/{device_id}/commands", response_model=schemas.CommandResult)
def send_device_command(device_id: int, command: schemas.DeviceCommand, db: Session = Depends(database.get_db)):
    """Отправить команду устройству"""
    device = db.query(models.Device).filter(
        models.Device.id == device_id,
        models.Device.is_active == True
    ).first()
    
    if not device:
        raise HTTPException(status_code=404, detail="Device not found")
    
    if device.status != "online":
        raise HTTPException(status_code=423, detail="Device is not available")
    
    # Simple command processing logic
    import uuid
    from datetime import datetime
    
    command_id = str(uuid.uuid4())
    
    # Simulate command execution
    result = schemas.CommandResult(
        command_id=command_id,
        status="accepted",
        message=f"Command '{command.command}' accepted for execution",
        timestamp=datetime.now()
    )
    
    # Store result (in production, use proper storage)
    command_results[command_id] = result
    
    return result

@router.get("/devices/{device_id}/state", response_model=schemas.DeviceState)
def get_device_state(device_id: int, db: Session = Depends(database.get_db)):
    """Получить текущее состояние устройства"""
    device = db.query(models.Device).filter(
        models.Device.id == device_id,
        models.Device.is_active == True
    ).first()
    
    if not device:
        raise HTTPException(status_code=404, detail="Device not found")
    
    # Simulate device state
    from datetime import datetime
    
    state = schemas.DeviceState(
        device_id=device_id,
        state={
            "status": device.status,
            "configuration": device.configuration,
            "last_seen": device.last_seen.isoformat() if device.last_seen else None
        },
        last_updated=datetime.now()
    )
    
    return state
