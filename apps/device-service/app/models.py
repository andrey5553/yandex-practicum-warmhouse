from sqlalchemy import Column, Integer, String, DateTime, JSON, Boolean
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.sql import func
import datetime

Base = declarative_base()

class Device(Base):
    __tablename__ = "devices"
    
    id = Column(Integer, primary_key=True, index=True)
    name = Column(String, nullable=False)
    type = Column(String, nullable=False)
    device_type_id = Column(Integer, nullable=False)
    room_id = Column(Integer)
    serial_number = Column(String, unique=True, nullable=False)
    status = Column(String, default="offline")
    configuration = Column(JSON, default=dict)
    last_seen = Column(DateTime, default=func.now())
    created_at = Column(DateTime, default=func.now())
    is_active = Column(Boolean, default=True)