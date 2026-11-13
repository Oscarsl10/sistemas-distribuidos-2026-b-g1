# Certification C3 -- Oscar Guillermo Sierra Lozano

## Cover

| Item | Detail |
|---|---|
| Full name | Oscar Guillermo Sierra Lozano |
| GitHub user | Oscarsl10 |
| Group (G1 or G2) | G1 |
| Team | CineSync Platform |
| Team -docs repository (full URL) | https://github.com/code-corhuila/csp-docs |
| Course (exact name) | HCCDP-Solution Architectures Course |
| Course link | https://e.huawei.com/en/talent/ict-academy/#/ict-courses-detail?courseId=6WtKiPd1S1YelBB0AUbE1vXlPgk |
| Start date (YYYY-MM-DD) | 2026-10-29 |
| Completion date (YYYY-MM-DD) | TODO |

## Study log

| Date (YYYY-MM-DD) | Hours | Chapters or labs | Evidence file |
|---|---|---|---|
| 2026-10-29 | 0.17 | Enrollment in the course (10 min) | evidence/00-enrollment.png |
| 2026-10-29 | 1.00 | Chapter 1 | evidence/01-chapter-01.png |
| 2026-10-30 | 4.00 | Chapter 2 | evidence/02-chapter-02.png |
| 2026-11-01 | 4.00 | Chapter 3 | evidence/03-chapter-03.png |
| 2026-11-03 | 4.00 | Chapter 4 | evidence/04-chapter-04.png |
| 2026-11-04 | 4.00 | Chapter 5 | evidence/05-chapter-05.png |
| 2026-11-05 | 4.00 | Chapter 6 | evidence/06-chapter-06.png |
| 2026-11-06 | 1.28 | Chapter 7 (part 1) | evidence/07-chapter-07.png |
| 2026-11-07 | 2.72 | Chapter 7 (part 2) | evidence/07-chapter-07.png |
| 2026-11-07 | 4.00 | Chapter 8 | evidence/08-chapter-08.png |
| 2026-11-08 | 4.00 | Chapter 9 | evidence/09-chapter-09.png |
| 2026-11-09 | 3.00 | Chapter 10 | evidence/10-chapter-10.png |
| 2026-11-11 | 3.00 | Chapter 11 | evidence/11-chapter-11.png |
| 2026-11-11 | 3.00 | Chapter 12 | evidence/12-chapter-12.png |
| 2026-11-12 | 3.00 | Chapter 13 | evidence/13-chapter-13.png |
| 2026-11-12 | 3.00 | Chapter 14 | evidence/14-chapter-14.png |
| 2026-11-13 | 2.07 | Lab 1 | evidence/15-lab-1.png |
| 2026-11-13 | 1.28 | Lab 2 | evidence/16-lab-2.png |
| 2026-11-13 | 1.77 | Lab 3 | evidence/17-lab-3.png |
| Total | 53.29 | | |

## Labs

| Evidence file | Lab name | What I built | What failed first and how I fixed it |
|---|---|---|---|
| evidence/15-lab-1.png | Design Cases and Practices of an O&M Solution | I created subnet-db and sg-db, bought the rds-web MySQL 5.7 primary/standby instance, and installed Discuz! against it. I created the SMN topic test with a confirmed email subscription and the Cloud Eye event alarm alarm-reboot on ecs-web. I installed ICAgent to ship /var/log/messages to LTS with the host group hccdp and a keyword alarm on Stopped. | ICAgent first showed Authentication error in LTS because the secret key in my install command was missing its first character. I compared the command with credentials.csv, fixed it, and the agent turned Running. The Discuz! page also timed out until I used http instead of https. The sandbox account could not create the LTS alarm notification rule, so the hccdp alarm has no notification. |
| evidence/16-lab-2.png | High-Performance Network Solution | I created subnet-db, sg-db and the rds-discuz MySQL 5.7 primary/standby instance, and installed Discuz! on ecs-web. I unbound the EIP, created the public NAT gateway nat-web with an SNAT rule for subnet-web and a DNAT rule for port 80, and checked that the forum still loaded. I created vpc-om with subnet-om, sg-om and the jump server ecs-om, connected the VPCs with the peering connection peering-web, and added a private DNS zone www.abc.com pointing to ecs-web. From ecs-om, both ping 192.168.1.161 and ping www.abc.com answered with 0% loss. | The Discuz! installer sent the form with localhost as the database server and failed to connect, so I went back and entered the RDS floating IP. The installer frame then stayed blank, and I confirmed the install by opening the forum and by finding install.lock on the server. I also reloaded the whole lab page by mistake, but the session and resources survived. |
| evidence/17-lab-3.png | Highly Elastic Service System | I created subnet-db, sg-db and the rds-web MySQL 5.7 primary/standby instance, installed Discuz! on ecs-web, unbound its EIP and put it behind the shared load balancer elb-web with listener-web on TCP 80. I built the image image-web, the AS configuration as-config-web and the AS group as-group-web (minimum 1, expected 2, maximum 5, attached to elb-web), with a CPU scale-out policy above 70% and a scale-in policy below 30%. I also created a bandwidth scaling policy that sets the elb-web EIP to 30 Mbit/s, added ecs-web to the group with instance protection, changed the AS configuration to c7n.large.4 and configured an SMN topic notification. | I could not type the underscore in the lab console, so the backend group is named server-group-web instead of server_group-web. My first AS configuration had no flavor selected and a 100 GiB disk, so I recreated it with s6.large.2 and 40 GiB. The CPU scale-out never triggered: I ran stress for 900 seconds and the CPU reached 100%, but the scale-in policy had already removed one instance and the alarm did not add one before the lab ended. |
| evidence/18-lab-4.png | TODO | TODO | TODO |

## Final exam

- Score: TODO
- Date (YYYY-MM-DD): TODO
- File: evidence/19-final-exam.png

## Skills summary

TODO (5 to 7 lines, English, first person, at least three concrete design techniques from the course)
