import type { components } from "./schema";

export type Schemas = components["schemas"];
export type Project = Schemas["Project"];
export type Member = Schemas["Member"];
export type APIKey = Schemas["APIKey"];
export type Provider = Schemas["Provider"];
export type Template = Schemas["Template"];
export type TemplateVersion = Schemas["TemplateVersion"];
export type Preview = Schemas["Preview"];
export type Message = Schemas["Message"];
export type MessageEvent = Schemas["MessageEvent"];
export type Batch = Schemas["Batch"];
export type Contact = Schemas["Contact"];
export type AdminDevice = Schemas["AdminDevice"];
export type WebhookDelivery = Schemas["WebhookDelivery"];
export type UsageRow = Schemas["UsageRow"];
export type Dashboard = Schemas["Dashboard"];
export type Health = Schemas["Health"];
export type AuditLog = Schemas["AuditLog"];
export type User = Schemas["User"];
export type APIError = Schemas["APIError"];
export type Channel = Schemas["Channel"];
export type MessageStatus = Schemas["MessageStatus"];
export type Role = Schemas["Role"];

export interface Envelope<T> {
  data: T;
  meta?: Record<string, unknown>;
  error?: APIError;
}

export type GatewayPairing = components["schemas"]["GatewayPairing"];
export type Group = Schemas["Group"];
export type AdminSendRequest = Schemas["AdminSendRequest"];
export type GroupMembersResult = Schemas["GroupMembersResult"];
export type InboundSMS = Schemas["InboundSMS"];
