import { check, sleep } from "k6";
import http from "k6/http";
import { vu } from 'k6/execution';
import { rand } from "./helper.js";
// TODO: replace with tiny version of faker
import { faker } from "https://esm.sh/@faker-js/faker";

export const options = {
  stages: [
    { duration: "5s", target: 40 },
    { duration: "10s", target: 90 },
    { duration: "5s", target: 15 },
  ],
};

export default function() {
  const url = "http://localhost:8080/api/v1";
  let res = http.get(`${url}/users/${rand(1,1000)}`);
  check(res, {
    "is status 200": (res) => res.status === 200,
    "is json a user struct": (res) => 
      res.json().hasOwnProperty("id") &&
      res.json().hasOwnProperty("email") &&
      res.json().hasOwnProperty("username") &&
      res.json().hasOwnProperty("fullname") &&
      res.json().hasOwnProperty("isMember") &&
      res.json().hasOwnProperty("internshipStartDate"),
  });

  const payload = {
    email: `${vu.idInTest}.${faker.internet.email()}`,
    username: `${vu.idInTest}.${faker.internet.username()}`,
    username: faker.person.fullName(),
    isMember: faker.datatype.boolean(),
    internshipStartDate: faker.date.anytime(),
  };
  res = http.post(`${url}/users`, JSON.stringify(payload), {
    headers: {"Content-Type": "application/json"},
  });
  check(res, {"is status 201": (r) => r.status === 201});
  sleep(1);
}
